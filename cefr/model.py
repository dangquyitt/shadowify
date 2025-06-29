import torch
import torch.nn as nn
import pandas as pd
import numpy as np
from torch.utils.data import Dataset, DataLoader
from transformers import BertTokenizer, BertModel, get_linear_schedule_with_warmup
from torch.optim import AdamW
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score
import matplotlib.pyplot as plt
import seaborn as sns
from tqdm import tqdm
import warnings
warnings.filterwarnings('ignore')

# Device setup
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
print(f"Using device: {device}")


class CEFRDataset(Dataset):
    """Dataset class for CEFR regression"""

    def __init__(self, texts, labels, tokenizer, max_length=512):
        self.texts = texts
        self.labels = labels
        self.tokenizer = tokenizer
        self.max_length = max_length

    def __len__(self):
        return len(self.texts)

    def __getitem__(self, idx):
        text = str(self.texts[idx])
        label = self.labels[idx]

        # Tokenize text
        encoding = self.tokenizer(
            text,
            truncation=True,
            padding='max_length',
            max_length=self.max_length,
            return_tensors='pt'
        )

        return {
            'input_ids': encoding['input_ids'].flatten(),
            'attention_mask': encoding['attention_mask'].flatten(),
            # Float for regression
            'labels': torch.tensor(label, dtype=torch.float)
        }


class BERTRegressor(nn.Module):
    """BERT model for regression"""

    def __init__(self, model_name='bert-base-uncased', dropout=0.3):
        super(BERTRegressor, self).__init__()
        self.bert = BertModel.from_pretrained(model_name)
        self.dropout = nn.Dropout(dropout)
        self.regressor = nn.Linear(self.bert.config.hidden_size, 1)

    def forward(self, input_ids, attention_mask):
        outputs = self.bert(input_ids=input_ids, attention_mask=attention_mask)
        pooled_output = outputs.pooler_output
        output = self.dropout(pooled_output)
        return self.regressor(output)


class CEFRRegressor:
    """BERT-based CEFR Level Regressor"""

    def __init__(self, model_name='bert-base-uncased'):
        self.model_name = model_name
        self.tokenizer = BertTokenizer.from_pretrained(model_name)
        self.model = BERTRegressor(model_name)
        self.model.to(device)

        # CEFR level mapping
        self.cefr_mapping = {
            1: 'A1', 2: 'A2', 3: 'B1',
            4: 'B2', 5: 'C1', 6: 'C2'
        }

    def load_data(self, train_path, val_path, test_path):
        """Load and preprocess data"""
        print("Loading data...")

        # Load CSV files
        train_df = pd.read_csv(train_path)
        val_df = pd.read_csv(val_path)
        test_df = pd.read_csv(test_path)

        # Use average of Annotator I and II as target value (1-6)
        train_df['target'] = (train_df['Annotator I'] +
                              train_df['Annotator II']) / 2
        val_df['target'] = (val_df['Annotator I'] + val_df['Annotator II']) / 2
        test_df['target'] = (test_df['Annotator I'] +
                             test_df['Annotator II']) / 2

        # Clip values to range 1-6
        train_df['target'] = train_df['target'].clip(1, 6)
        val_df['target'] = val_df['target'].clip(1, 6)
        test_df['target'] = test_df['target'].clip(1, 6)

        print(f"Train samples: {len(train_df)}")
        print(f"Validation samples: {len(val_df)}")
        print(f"Test samples: {len(test_df)}")

        # Display target distribution
        print("\nTarget distribution in training set:")
        print(train_df['target'].describe())

        return train_df, val_df, test_df

    def create_data_loaders(self, train_df, val_df, test_df, batch_size=16, max_length=512):
        """Create DataLoaders"""
        print("Creating data loaders...")

        # Create datasets
        train_dataset = CEFRDataset(
            train_df['text'].values,
            train_df['target'].values,  # Use target instead of label
            self.tokenizer,
            max_length
        )

        val_dataset = CEFRDataset(
            val_df['text'].values,
            val_df['target'].values,
            self.tokenizer,
            max_length
        )

        test_dataset = CEFRDataset(
            test_df['text'].values,
            test_df['target'].values,
            self.tokenizer,
            max_length
        )

        # Create data loaders
        train_loader = DataLoader(
            train_dataset, batch_size=batch_size, shuffle=True)
        val_loader = DataLoader(
            val_dataset, batch_size=batch_size, shuffle=False)
        test_loader = DataLoader(
            test_dataset, batch_size=batch_size, shuffle=False)

        return train_loader, val_loader, test_loader

    def train(self, train_loader, val_loader, epochs=3, learning_rate=2e-5):
        """Training function for regression"""
        print(f"Starting training for {epochs} epochs...")

        # Loss function for regression
        criterion = nn.MSELoss()

        # Optimizer and scheduler
        optimizer = AdamW(self.model.parameters(), lr=learning_rate)
        total_steps = len(train_loader) * epochs
        scheduler = get_linear_schedule_with_warmup(
            optimizer,
            num_warmup_steps=0,
            num_training_steps=total_steps
        )

        train_losses = []
        val_mse_scores = []

        for epoch in range(epochs):
            print(f'\nEpoch {epoch + 1}/{epochs}')

            # Training
            self.model.train()
            total_train_loss = 0

            train_bar = tqdm(train_loader, desc='Training')
            for batch in train_bar:
                optimizer.zero_grad()

                input_ids = batch['input_ids'].to(device)
                attention_mask = batch['attention_mask'].to(device)
                labels = batch['labels'].to(device)

                outputs = self.model(input_ids=input_ids,
                                     attention_mask=attention_mask)
                loss = criterion(outputs.squeeze(), labels)

                total_train_loss += loss.item()

                loss.backward()
                torch.nn.utils.clip_grad_norm_(self.model.parameters(), 1.0)
                optimizer.step()
                scheduler.step()

                train_bar.set_postfix({'loss': loss.item()})

            avg_train_loss = total_train_loss / len(train_loader)
            train_losses.append(avg_train_loss)

            # Validation
            val_mse = self.evaluate(val_loader)
            val_mse_scores.append(val_mse)

            print(f'Average training loss: {avg_train_loss:.4f}')
            print(f'Validation MSE: {val_mse:.4f}')
            print(f'Validation RMSE: {np.sqrt(val_mse):.4f}')

        return train_losses, val_mse_scores

    def evaluate(self, data_loader):
        """Evaluation function for regression"""
        self.model.eval()
        predictions = []
        true_values = []

        with torch.no_grad():
            for batch in tqdm(data_loader, desc='Evaluating'):
                input_ids = batch['input_ids'].to(device)
                attention_mask = batch['attention_mask'].to(device)
                labels = batch['labels'].to(device)

                outputs = self.model(input_ids=input_ids,
                                     attention_mask=attention_mask)

                predictions.extend(outputs.squeeze().cpu().tolist())
                true_values.extend(labels.cpu().tolist())

        mse = mean_squared_error(true_values, predictions)
        return mse

    def predict(self, text, max_length=512):
        """Predict CEFR level for a sentence (regression)"""
        self.model.eval()

        encoding = self.tokenizer(
            text,
            truncation=True,
            padding='max_length',
            max_length=max_length,
            return_tensors='pt'
        )

        input_ids = encoding['input_ids'].to(device)
        attention_mask = encoding['attention_mask'].to(device)

        with torch.no_grad():
            output = self.model(input_ids=input_ids,
                                attention_mask=attention_mask)
            predicted_value = output.squeeze().cpu().item()

        # Clip value to range 1-6
        predicted_value = max(1.0, min(6.0, predicted_value))

        # Map to CEFR level
        rounded_level = round(predicted_value)
        cefr_level = self.cefr_mapping.get(rounded_level, 'Unknown')

        return {
            'predicted_value': predicted_value,
            'rounded_level': rounded_level,
            'cefr_level': cefr_level
        }

    def detailed_evaluation(self, data_loader):
        """Detailed evaluation with regression metrics"""
        self.model.eval()
        predictions = []
        true_values = []

        with torch.no_grad():
            for batch in tqdm(data_loader, desc='Detailed Evaluation'):
                input_ids = batch['input_ids'].to(device)
                attention_mask = batch['attention_mask'].to(device)
                labels = batch['labels'].to(device)

                outputs = self.model(input_ids=input_ids,
                                     attention_mask=attention_mask)

                predictions.extend(outputs.squeeze().cpu().tolist())
                true_values.extend(labels.cpu().tolist())

        # Regression metrics
        mse = mean_squared_error(true_values, predictions)
        rmse = np.sqrt(mse)
        mae = mean_absolute_error(true_values, predictions)
        r2 = r2_score(true_values, predictions)

        print("Regression Metrics:")
        print(f"MSE: {mse:.4f}")
        print(f"RMSE: {rmse:.4f}")
        print(f"MAE: {mae:.4f}")
        print(f"R²: {r2:.4f}")

        # Scatter plot
        plt.figure(figsize=(10, 8))
        plt.scatter(true_values, predictions, alpha=0.6)
        plt.plot([1, 6], [1, 6], 'r--', lw=2)
        plt.xlabel('True CEFR Level')
        plt.ylabel('Predicted CEFR Level')
        plt.title('True vs Predicted CEFR Levels')
        plt.xlim(0.5, 6.5)
        plt.ylim(0.5, 6.5)
        plt.grid(True, alpha=0.3)

        # Add R² to plot
        plt.text(0.05, 0.95, f'R² = {r2:.3f}', transform=plt.gca().transAxes,
                 bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.5))

        plt.tight_layout()
        plt.savefig('regression_plot.png', dpi=300, bbox_inches='tight')
        plt.show()

        # Rounded predictions for classification-like evaluation
        rounded_predictions = [max(1, min(6, round(p))) for p in predictions]
        rounded_true = [round(t) for t in true_values]

        # Classification accuracy on rounded values
        accuracy = sum(rp == rt for rp, rt in zip(
            rounded_predictions, rounded_true)) / len(rounded_predictions)
        print(f"\nClassification accuracy (rounded): {accuracy:.4f}")

        return mse, predictions, true_values

    def save_model(self, path='./cefr_bert_regressor'):
        """Save model"""
        torch.save({
            'model_state_dict': self.model.state_dict(),
            'model_name': self.model_name
        }, f"{path}/model.pt")
        self.tokenizer.save_pretrained(path)
        print(f"Model saved to {path}")

    def load_model(self, path='./cefr_bert_regressor'):
        """Load model"""
        checkpoint = torch.load(f"{path}/model.pt", map_location=device)
        self.model = BERTRegressor(checkpoint['model_name'])
        self.model.load_state_dict(checkpoint['model_state_dict'])
        self.model.to(device)
        self.tokenizer = BertTokenizer.from_pretrained(path)
        print(f"Model loaded from {path}")


def main():
    """Main function to train and test model"""
    # Initialize regressor
    regressor = CEFRRegressor()

    # Load data
    train_df, val_df, test_df = regressor.load_data(
        'dataset/train.csv',
        'dataset/validation.csv',
        'dataset/test.csv'
    )

    # Create data loaders
    train_loader, val_loader, test_loader = regressor.create_data_loaders(
        train_df, val_df, test_df, batch_size=16
    )

    # Training
    train_losses, val_mse_scores = regressor.train(
        train_loader, val_loader, epochs=3, learning_rate=2e-5
    )

    # Evaluation on test set
    print("\n" + "="*50)
    print("FINAL EVALUATION ON TEST SET")
    print("="*50)

    test_mse, predictions, true_values = regressor.detailed_evaluation(
        test_loader)
    print(f"\nFinal Test MSE: {test_mse:.4f}")
    print(f"Final Test RMSE: {np.sqrt(test_mse):.4f}")

    # Save model
    regressor.save_model()

    # Test with sample sentences
    print("\n" + "="*50)
    print("SAMPLE PREDICTIONS")
    print("="*50)

    sample_texts = [
        "I like cats.",  # Simple sentence - A1
        "She had a beautiful necklace around her neck.",  # A1/A2
        "The weather is getting warmer as spring approaches.",  # B1/B2
        "Renewable energy projects in many developing countries have demonstrated that renewable energy can directly contribute to poverty alleviation.",  # C1/C2
    ]

    for text in sample_texts:
        result = regressor.predict(text)
        print(f"Text: {text}")
        print(f"Predicted Value: {result['predicted_value']:.2f}")
        print(
            f"Predicted CEFR Level: {result['cefr_level']} (Level {result['rounded_level']})")
        print("-" * 50)


if __name__ == "__main__":
    main()
