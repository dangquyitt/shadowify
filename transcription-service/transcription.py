import whisper
import json
import yt_dlp
import os
import tempfile


def transcribe_audio(url):
    # Tạo file tạm để tránh ghi đè
    with tempfile.NamedTemporaryFile(suffix=".m4a") as temp_audio:
        output_file = temp_audio.name

    ydl_opts = {
        'format': 'm4a/bestaudio/best',
        'postprocessors': [{'key': 'FFmpegExtractAudio', 'preferredcodec': 'm4a'}],
        'outtmpl': output_file
    }

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        error_code = ydl.download([url])
        if error_code:
            raise RuntimeError("Download failed.")

    # Đảm bảo lấy đúng file đầu ra
    if not os.path.exists(output_file):
        output_file += ".m4a"

    # Load model Whisper
    model = whisper.load_model("tiny")  # Thay bằng model khác nếu cần
    print("Whisper model loaded.")

    # Nhận diện giọng nói
    transcribe = model.transcribe(audio=output_file)
    segments = transcribe.get('segments', [])

    # Chuyển đổi kết quả thành JSON
    transcripts = [
        {"id": seg['id'] + 1, "startTime": int(seg['start']), "endTime": int(
            seg['end']), "text": seg['text'].strip()}
        for seg in segments
    ]

    transcript_json = json.dumps(transcripts, indent=2, ensure_ascii=False)

    # Xóa file tạm sau khi hoàn thành
    try:
        os.remove(output_file)
    except OSError:
        pass

    print("Transcription complete.")

    return transcript_json
