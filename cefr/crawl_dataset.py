import pandas as pd

splits = {
    'train': 'data/train-00000-of-00001.parquet',
    'validation': 'data/validation-00000-of-00001.parquet',
    'test': 'data/test-00000-of-00001.parquet'
}

for split_name, relative_path in splits.items():
    parquet_path = "hf://datasets/edesaras/CEFR-Sentence-Level-Annotations/" + relative_path
    print(f"📥 Reading {split_name} từ {parquet_path}")

    df = pd.read_parquet(parquet_path, engine='pyarrow')

    output_file = f"cefr_{split_name}.csv"
    df.to_csv(output_file, index=False, encoding='utf-8')

    print(f"✅ Write {output_file} with {len(df)} lines.")
