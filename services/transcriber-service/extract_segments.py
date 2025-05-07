from faster_whisper import WhisperModel


def extract_segments(file_path: str, model_size: str = "distil-large-v3") -> list:
    model = WhisperModel(model_size_or_path=model_size)

    segments_list = []
    segments, info = model.transcribe(
        file_path,
        beam_size=5,
        log_progress=True
    )

    print("Detected language '%s' with probability %f" %
          (info.language, info.language_probability))

    for segment in segments:
        segments_list.append({
            "start": round(segment.start, 2),
            "end": round(segment.end, 2),
            "text": segment.text.strip()
        })

    return segments_list


# Example usage
if __name__ == "__main__":
    audio_path = "./downloads/9QmEaAUSjm8_20250508_000758.mp3"
    result = extract_segments(audio_path)

    for seg in result:
        print(f"[{seg['start']}s -> {seg['end']}s] {seg['text']}")
