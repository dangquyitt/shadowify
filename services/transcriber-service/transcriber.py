from faster_whisper import WhisperModel


class Transcriber:
    def __init__(self, model_size="tiny", device="auto", compute_type="int8"):
        self.model = WhisperModel(
            model_size_or_path=model_size,
            device=device,
            compute_type=compute_type
        )

    def extract_segments(self, file_path: str) -> list:
        segments, info = self.model.transcribe(
            file_path,
            beam_size=5,
            log_progress=True
        )

        print(
            f"🔤 Detected language: {info.language} ({info.language_probability:.2f})")

        return [
            {
                "start": round(segment.start, 2),
                "end": round(segment.end, 2),
                "text": segment.text.strip()
            }
            for segment in segments
        ]
