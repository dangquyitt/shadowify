import pandas as pd
import os

# Create dataset directory if it doesn't exist
dataset_dir = "dataset"
if not os.path.exists(dataset_dir):
    os.makedirs(dataset_dir)
    print(f"📁 Created directory: {dataset_dir}")

splits = {
    'train': 'data/train-00000-of-00001.parquet',
    'validation': 'data/validation-00000-of-00001.parquet',
    'test': 'data/test-00000-of-00001.parquet'
}

for split_name, relative_path in splits.items():
    parquet_path = "hf://datasets/edesaras/CEFR-Sentence-Level-Annotations/" + relative_path
    print(f"📥 Reading {split_name} from {parquet_path}")

    df = pd.read_parquet(parquet_path, engine='pyarrow')

    output_file = os.path.join(dataset_dir, f"{split_name}.csv")
    df.to_csv(output_file, index=False, encoding='utf-8')

    print(f"✅ Saved {output_file} with {len(df)} rows.")
    print(f"   Columns: {list(df.columns)}")
    print(f"   Shape: {df.shape}")
    print()

print("🎉 Dataset download completed!")
print(f"📂 All files saved in: {os.path.abspath(dataset_dir)}")
