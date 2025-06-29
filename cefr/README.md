# CEFR Level Prediction using BERT Regression

This project implements a BERT-based regression model to predict CEFR (Common European Framework of Reference for Languages) levels from English sentences.

## 📋 Overview

The model predicts CEFR levels as continuous values from 1-6, which are then mapped to discrete levels:

- **A1** (Beginner): Can understand and use familiar everyday expressions
- **A2** (Elementary): Can communicate in simple tasks requiring simple exchange of information
- **B1** (Intermediate): Can deal with most situations likely to arise whilst travelling
- **B2** (Upper Intermediate): Can interact with fluency and spontaneity
- **C1** (Advanced): Can express ideas fluently and spontaneously
- **C2** (Proficient): Can understand virtually everything heard or read

## 🗂️ Dataset Structure

The dataset contains CSV files with the following structure:

```
text,Annotator I,Annotator II
"Example sentence",1,1
```

Where:

- `text`: The English sentence to classify
- `Annotator I` & `Annotator II`: CEFR level ratings (1-6) from two annotators
- Levels 1-6 correspond to A1-C2 respectively

## 🚀 Quick Start

### 1. Setup Environment

```bash
# Create virtual environment
python -m venv .venv
source .venv/bin/activate  # On macOS/Linux
# or
.venv\Scripts\activate     # On Windows

# Install dependencies
pip install -r requirements.txt
```

### 2. Train the Model

```bash
python model.py
```

This will:

- Load and preprocess the dataset
- Train a BERT model for 3 epochs
- Evaluate on test set
- Save the trained model
- Generate classification report and confusion matrix

### 3. Demo Predictions

```bash
# Run demo with predefined sentences
python demo.py

# Run interactive demo
python demo.py --interactive
```

## 🧠 Model Architecture

- **Base Model**: `bert-base-uncased`
- **Task**: Multi-class classification (6 classes)
- **Input**: English text (max 512 tokens)
- **Output**: CEFR level (A1-C2) with confidence score

## 📊 Model Performance

The model uses:

- **Loss Function**: Cross-entropy loss
- **Optimizer**: AdamW with learning rate 2e-5
- **Scheduler**: Linear warmup scheduler
- **Batch Size**: 16
- **Max Length**: 512 tokens
- **Training Epochs**: 3

## 🔧 Usage Examples

### Basic Prediction

```python
from model import CEFRClassifier

# Initialize classifier
classifier = CEFRClassifier()

# Load trained model
classifier.load_model('./cefr_bert_model')

# Predict CEFR level
result = classifier.predict("I like cats.")
print(f"CEFR Level: {result['cefr_level']}")
print(f"Confidence: {result['confidence']:.2%}")
```

### Training Custom Model

```python
from model import CEFRClassifier

# Initialize classifier
classifier = CEFRClassifier()

# Load your data
train_df, val_df, test_df = classifier.load_data(
    'path/to/train.csv',
    'path/to/validation.csv',
    'path/to/test.csv'
)

# Create data loaders
train_loader, val_loader, test_loader = classifier.create_data_loaders(
    train_df, val_df, test_df, batch_size=16
)

# Train model
classifier.train(train_loader, val_loader, epochs=3)

# Evaluate
accuracy = classifier.detailed_evaluation(test_loader)
print(f"Test Accuracy: {accuracy:.4f}")

# Save model
classifier.save_model('./my_cefr_model')
```

## 📁 File Structure

```
cefr/
├── model.py              # Main BERT model implementation
├── demo.py               # Demo and interactive prediction
├── requirements.txt      # Python dependencies
├── README.md            # This file
├── dataset/             # Training data
│   ├── train.csv
│   ├── validation.csv
│   └── test.csv
└── cefr_bert_model/     # Saved model (after training)
    ├── config.json
    ├── pytorch_model.bin
    └── tokenizer files
```

## 🎯 Features

- **BERT-based Classification**: Uses pre-trained BERT for robust text understanding
- **Multi-annotator Support**: Averages ratings from multiple annotators
- **Detailed Evaluation**: Provides classification report and confusion matrix
- **Interactive Demo**: Test model with custom sentences
- **Model Persistence**: Save/load trained models
- **Confidence Scores**: Get prediction confidence for each classification

## 📈 Output Examples

```
Text: "I like cats."
Predicted CEFR Level: A1 (Level 1)
Confidence: 92.5%

Text: "The government should implement more effective policies."
Predicted CEFR Level: B2 (Level 4)
Confidence: 87.3%
```

## 🔧 Requirements

- Python 3.7+
- PyTorch 2.0+
- Transformers 4.20+
- scikit-learn 1.0+
- pandas 1.3+
- Other dependencies in `requirements.txt`

## 💡 Tips

1. **GPU Training**: The model will automatically use GPU if available for faster training
2. **Batch Size**: Adjust batch size based on your GPU memory
3. **Fine-tuning**: Experiment with learning rates and epochs for better performance
4. **Data Quality**: Ensure consistent annotation quality for better results

## 🤝 Contributing

Feel free to submit issues and enhancement requests!

## 📄 License

This project is for educational purposes.
