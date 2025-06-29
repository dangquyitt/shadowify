# CEFR Level Classification using BERT

This project implements a BERT-based text classifier to predict CEFR (Common European Framework of Reference for Languages) levels from English sentences.

## 📋 Overview

The model classifies English text into 6 CEFR levels using a **classification approach**:

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
- The model uses the average of both annotators, rounded to nearest integer, then converted to class labels (0-5)

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

- Load and preprocess the dataset (convert ratings to class labels)
- Train a BERT classification model for 3 epochs
- Evaluate on test set with accuracy, classification report, and confusion matrix
- Save the trained model
- Test with sample predictions

## 🧠 Model Architecture

- **Base Model**: `bert-base-uncased`
- **Task**: Multi-class classification (6 classes: A1-C2)
- **Input**: English text (max 512 tokens)
- **Output**: CEFR level (A1-C2) with confidence score
- **Classes**: 0-5 (mapped to A1-C2)

## 📊 Model Performance

The model uses:

- **Loss Function**: Cross-entropy loss (built into BertForSequenceClassification)
- **Optimizer**: AdamW with learning rate 2e-5
- **Scheduler**: Linear warmup scheduler
- **Batch Size**: 16
- **Max Length**: 512 tokens
- **Training Epochs**: 3
- **Evaluation Metrics**: Accuracy, Classification Report, Confusion Matrix

## 🔧 Usage Examples

### Basic Prediction

```python
from model import CEFRClassifier

# Initialize classifier
classifier = CEFRClassifier()

# Load trained model
classifier.load_model('./cefr_bert_classifier')

# Predict CEFR level
# Predict CEFR level
result = classifier.predict("I like cats.")
print(f"CEFR Level: {result['cefr_level']}")
print(f"Confidence: {result['confidence']:.2%}")
print(f"Predicted Class: {result['predicted_class']}")
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
train_losses, val_accuracies = classifier.train(
    train_loader, val_loader, epochs=3
)

# Detailed evaluation
accuracy, predictions, true_labels = classifier.detailed_evaluation(test_loader)
print(f"Test Accuracy: {accuracy:.4f}")

# Save model
classifier.save_model('./my_cefr_model')
```

## 📁 File Structure

```
cefr/
├── model.py              # Main BERT classification implementation
├── requirements.txt      # Python dependencies
├── README.md            # This file
├── crawl_dataset.py     # Dataset crawling script
├── dataset/             # Training data
│   ├── train.csv
│   ├── validation.csv
│   └── test.csv
└── cefr_bert_classifier/ # Saved model (after training)
    ├── config.json
    ├── pytorch_model.bin
    └── tokenizer files
```

## 🎯 Features

- **BERT Classification**: Uses pre-trained BERT with classification head
- **Multi-annotator Support**: Averages ratings from multiple annotators
- **Robust Preprocessing**: Converts continuous ratings to discrete classes
- **Comprehensive Evaluation**: Classification report, confusion matrix, accuracy
- **Confidence Scores**: Get prediction confidence for each classification
- **Model Persistence**: Save/load trained models using HuggingFace format

## 📈 Output Examples

```
Text: "I like cats."
Predicted CEFR Level: A1 (Level 1)
Confidence: 92.5%

Text: "The government should implement more effective policies."
Predicted CEFR Level: B2 (Level 4)
Confidence: 87.3%
```

## 🔧 Dependencies

```
torch>=2.0.0
transformers>=4.40.0
scikit-learn>=1.4.0
pandas>=2.2.0
numpy>=1.26.0
tqdm>=4.66.0
matplotlib>=3.8.0
seaborn>=0.13.0
```

## 🆚 Classification vs Regression

This implementation uses **classification approach** which offers:

### ✅ **Advantages:**

- **Interpretable confidence scores** for each prediction
- **Robust to outliers** - no predictions outside valid range
- **Standard metrics** - accuracy, precision, recall, F1-score
- **Better for discrete levels** - CEFR levels are inherently categorical

### 📊 **Evaluation Metrics:**

- **Accuracy**: Overall classification accuracy
- **Classification Report**: Per-class precision, recall, F1-score
- **Confusion Matrix**: Visual representation of prediction vs actual

## 💡 Tips

1. **GPU Training**: Model automatically uses GPU if available for faster training
2. **Batch Size**: Adjust based on your GPU memory (16 works well for most setups)
3. **Epochs**: 3 epochs usually sufficient, monitor validation accuracy
4. **Data Quality**: Ensure consistent annotation quality between annotators
5. **Class Balance**: Check label distribution - model performs better with balanced classes

## 🤝 Contributing

Feel free to submit issues and enhancement requests!

## 📄 License

This project is for educational and research purposes.

## 📄 License

This project is for educational purposes.
