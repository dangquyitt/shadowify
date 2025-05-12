from yt_dlp import YoutubeDL
import os


class AudioDownloader:
    def __init__(self, output_dir="/tmp", format="wav"):
        self.output_dir = output_dir
        self.format = format
        os.makedirs(self.output_dir, exist_ok=True)

        self.ydl_opts = {
            'format': 'bestaudio/best',
            'outtmpl': os.path.join(self.output_dir, '%(id)s_%(upload_date)s.%(ext)s'),
            'noplaylist': True,
            'quiet': True,
            'postprocessors': [{
                'key': 'FFmpegExtractAudio',
                'preferredcodec': format,
                'preferredquality': '192',
            }],
            'prefer_ffmpeg': True,
        }
        self.ydl = YoutubeDL(self.ydl_opts)

    def download(self, video_code: str) -> str:
        url = f'https://www.youtube.com/watch?v={video_code}'
        try:
            info = self.ydl.extract_info(url, download=True)
            filename = os.path.splitext(self.ydl.prepare_filename(info))[
                0] + f".{self.format}"
            if os.path.exists(filename):
                return filename
            else:
                raise FileNotFoundError(
                    f"Downloaded file not found: {filename}")
        except Exception as e:
            raise RuntimeError(f"❌ Download failed for {url}: {e}")
