import os
import yt_dlp
from datetime import datetime

URLS = ['https://www.youtube.com/watch?v=']


def download_audio_from_youtube(video_code: str, output_dir: str = './downloads', format: str = 'mp3') -> str:
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    url = f'https://www.youtube.com/watch?v={video_code}'
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    filename = f"{video_code}_{timestamp}"

    ydl_opts = {
        'format': 'bestaudio/best',
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': format,
        }],
        'outtmpl': os.path.join(output_dir, f'{filename}.%(ext)s'),
        'noplaylist': True,
        'quiet': True,
    }

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        error_code = ydl.download(url)

    if error_code == 0:
        audio_file = os.path.join(output_dir, f'{filename}.{format}')
        return audio_file
    else:
        raise Exception(f"Download failed for URL: {url}")


# Ví dụ chạy thử
if __name__ == '__main__':
    video_code = '9QmEaAUSjm8'  # demo video của yt-dlp
    path = download_audio_from_youtube(video_code)
    print(f'✅ Audio saved at: {path}')
