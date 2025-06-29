import torch
import pandas as pd
import numpy as np
from torch.utils.data import Dataset, DataLoader
from transformers import BertTokenizer, BertForSequenceClassification, get_linear_schedule_with_warmup
from torch.optim import AdamW
from sklearn.metrics import accuracy_score, classification_report, confusion_matrix
import matplotlib.pyplot as plt
import seaborn as sns
from tqdm import tqdm
import warnings
warnings.filterwarnings('ignore')

# Device setup
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
print(f"Using device: {device}")


class CEFRDataset(Dataset):
    """Dataset class for CEFR classification"""

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
            # Long for classification
            'labels': torch.tensor(label, dtype=torch.long)
        }


class CEFRClassifier:
    """BERT-based CEFR Level Classifier"""

    def __init__(self, num_classes=6, model_name='bert-base-uncased'):
        self.num_classes = num_classes
        self.model_name = model_name
        self.tokenizer = BertTokenizer.from_pretrained(model_name)
        self.model = BertForSequenceClassification.from_pretrained(
            model_name,
            num_labels=num_classes
        )
        self.model.to(device)

        # CEFR level mapping (0-5 to A1-C2)
        self.cefr_mapping = {
            0: 'A1', 1: 'A2', 2: 'B1',
            3: 'B2', 4: 'C1', 5: 'C2'
        }

    def load_data(self, train_path, val_path, test_path):
        """Load and preprocess data"""
        print("Loading data...")

        # Load CSV files
        train_df = pd.read_csv(train_path)
        val_df = pd.read_csv(val_path)
        test_df = pd.read_csv(test_path)

        # Use average of Annotator I and II as target value, then convert to class (0-5)
        train_df['target'] = (train_df['Annotator I'] +
                              train_df['Annotator II']) / 2
        val_df['target'] = (val_df['Annotator I'] + val_df['Annotator II']) / 2
        test_df['target'] = (test_df['Annotator I'] +
                             test_df['Annotator II']) / 2

        # Round and convert to class labels (1-6 -> 0-5)
        train_df['label'] = (train_df['target'].round() - 1).astype(int)
        val_df['label'] = (val_df['target'].round() - 1).astype(int)
        test_df['label'] = (test_df['target'].round() - 1).astype(int)

        # Clip values to range 0-5
        train_df['label'] = train_df['label'].clip(0, 5)
        val_df['label'] = val_df['label'].clip(0, 5)
        test_df['label'] = test_df['label'].clip(0, 5)

        print(f"Train samples: {len(train_df)}")
        print(f"Validation samples: {len(val_df)}")
        print(f"Test samples: {len(test_df)}")

        # Display label distribution
        print("\nLabel distribution in training set:")
        print(train_df['label'].value_counts().sort_index())

        return train_df, val_df, test_df

    def create_data_loaders(self, train_df, val_df, test_df, batch_size=16, max_length=512):
        """Create DataLoaders"""
        print("Creating data loaders...")

        # Create datasets
        train_dataset = CEFRDataset(
            train_df['text'].values,
            train_df['label'].values,  # Use label for classification
            self.tokenizer,
            max_length
        )

        val_dataset = CEFRDataset(
            val_df['text'].values,
            val_df['label'].values,
            self.tokenizer,
            max_length
        )

        test_dataset = CEFRDataset(
            test_df['text'].values,
            test_df['label'].values,
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
        """Training function for classification"""
        print(f"Starting training for {epochs} epochs...")

        # Optimizer and scheduler
        optimizer = AdamW(self.model.parameters(), lr=learning_rate)
        total_steps = len(train_loader) * epochs
        scheduler = get_linear_schedule_with_warmup(
            optimizer,
            num_warmup_steps=0,
            num_training_steps=total_steps
        )

        train_losses = []
        val_accuracies = []

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

                outputs = self.model(
                    input_ids=input_ids,
                    attention_mask=attention_mask,
                    labels=labels
                )

                loss = outputs.loss
                total_train_loss += loss.item()

                loss.backward()
                torch.nn.utils.clip_grad_norm_(self.model.parameters(), 1.0)
                optimizer.step()
                scheduler.step()

                train_bar.set_postfix({'loss': loss.item()})

            avg_train_loss = total_train_loss / len(train_loader)
            train_losses.append(avg_train_loss)

            # Validation
            val_accuracy = self.evaluate(val_loader)
            val_accuracies.append(val_accuracy)

            print(f'Average training loss: {avg_train_loss:.4f}')
            print(f'Validation accuracy: {val_accuracy:.4f}')

        return train_losses, val_accuracies

    def evaluate(self, data_loader):
        """Evaluation function for classification"""
        self.model.eval()
        predictions = []
        true_labels = []

        with torch.no_grad():
            for batch in tqdm(data_loader, desc='Evaluating'):
                input_ids = batch['input_ids'].to(device)
                attention_mask = batch['attention_mask'].to(device)
                labels = batch['labels'].to(device)

                outputs = self.model(
                    input_ids=input_ids,
                    attention_mask=attention_mask
                )

                _, preds = torch.max(outputs.logits, dim=1)
                predictions.extend(preds.cpu().tolist())
                true_labels.extend(labels.cpu().tolist())

        accuracy = accuracy_score(true_labels, predictions)
        return accuracy

    def predict(self, text, max_length=512):
        """Predict CEFR level for a sentence (classification)"""
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
            outputs = self.model(input_ids=input_ids,
                                 attention_mask=attention_mask)
            _, prediction = torch.max(outputs.logits, dim=1)

        predicted_class = prediction.cpu().item()
        cefr_level = self.cefr_mapping[predicted_class]
        confidence = torch.softmax(outputs.logits, dim=1).max().cpu().item()

        return {
            'predicted_class': predicted_class + 1,  # Convert back to 1-6
            'cefr_level': cefr_level,
            'confidence': confidence
        }

    def detailed_evaluation(self, data_loader):
        """Detailed evaluation with classification metrics"""
        self.model.eval()
        predictions = []
        true_labels = []

        with torch.no_grad():
            for batch in tqdm(data_loader, desc='Detailed Evaluation'):
                input_ids = batch['input_ids'].to(device)
                attention_mask = batch['attention_mask'].to(device)
                labels = batch['labels'].to(device)

                outputs = self.model(
                    input_ids=input_ids,
                    attention_mask=attention_mask
                )

                _, preds = torch.max(outputs.logits, dim=1)
                predictions.extend(preds.cpu().tolist())
                true_labels.extend(labels.cpu().tolist())

        # Classification report
        target_names = [
            f'{self.cefr_mapping[i]} (Level {i+1})' for i in range(6)]
        report = classification_report(
            true_labels, predictions, target_names=target_names)
        print("Classification Report:")
        print(report)

        # Confusion Matrix
        cm = confusion_matrix(true_labels, predictions)
        plt.figure(figsize=(10, 8))
        sns.heatmap(cm, annot=True, fmt='d', cmap='Blues',
                    xticklabels=target_names, yticklabels=target_names)
        plt.title('Confusion Matrix - CEFR Level Classification')
        plt.ylabel('True Label')
        plt.xlabel('Predicted Label')
        plt.xticks(rotation=45)
        plt.yticks(rotation=0)
        plt.tight_layout()
        plt.savefig('confusion_matrix.png', dpi=300, bbox_inches='tight')
        plt.show()

        accuracy = accuracy_score(true_labels, predictions)
        return accuracy, predictions, true_labels

    def save_model(self, path='./cefr_bert_classifier'):
        """Save model"""
        self.model.save_pretrained(path)
        self.tokenizer.save_pretrained(path)
        print(f"Model saved to {path}")

    def load_model(self, path='./cefr_bert_classifier'):
        """Load model"""
        self.model = BertForSequenceClassification.from_pretrained(path)
        self.tokenizer = BertTokenizer.from_pretrained(path)
        self.model.to(device)
        print(f"Model loaded from {path}")


def main():
    """Main function to train and test model"""
    # Initialize classifier
    classifier = CEFRClassifier()

    # Load data
    train_df, val_df, test_df = classifier.load_data(
        'dataset/train.csv',
        'dataset/validation.csv',
        'dataset/test.csv'
    )

    # Create data loaders
    train_loader, val_loader, test_loader = classifier.create_data_loaders(
        train_df, val_df, test_df, batch_size=16
    )

    # Training
    train_losses, val_accuracies = classifier.train(
        train_loader, val_loader, epochs=3, learning_rate=2e-5
    )

    # Evaluation on test set
    print("\n" + "="*50)
    print("FINAL EVALUATION ON TEST SET")
    print("="*50)

    test_accuracy, predictions, true_labels = classifier.detailed_evaluation(
        test_loader)
    print(f"\nFinal Test Accuracy: {test_accuracy:.4f}")

    # Save model
    classifier.save_model()

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
        result = classifier.predict(text)
        print(f"Text: {text}")
        print(
            f"Predicted CEFR Level: {result['cefr_level']} (Level {result['predicted_class']})")
        print(f"Confidence: {result['confidence']:.4f}")
        print("-" * 50)


if __name__ == "__main__":
    main()
