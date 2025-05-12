# metadata_extractor.py
import yt_dlp


class MetadataExtractor:
    def extract(self, video_code: str) -> dict:
        url = f"https://www.youtube.com/watch?v={video_code}"

        ydl_opts = {
            'quiet': True,
            'skip_download': True,
            'simulate': True,
        }

        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            return ydl.extract_info(url, download=False)
