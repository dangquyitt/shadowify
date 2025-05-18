#!/usr/bin/env python3
import os
import json
import logging
import tempfile
import subprocess
from pathlib import Path
import yt_dlp
from faster_whisper import WhisperModel

logger = logging.getLogger(__name__)


class YoutubeExtractor:
    def __init__(self, config):
        self.config = config
        self.temp_dir = config.get('temp_dir', '/tmp')
        self.whisper_model_size = config.get('whisper_model_size', 'base')
        self.whisper_device = config.get('whisper_device', 'cpu')
        self.whisper_compute_type = config.get(
            'whisper_compute_type', 'float32')

        # Initialize Whisper model
        logger.info(f"Initializing Whisper model: {self.whisper_model_size}")
        self.whisper_model = WhisperModel(
            self.whisper_model_size,
            device=self.whisper_device,
            compute_type=self.whisper_compute_type
        )
        logger.info("Whisper model initialized")

    def extract_metadata(self, youtube_id):
        """
        Extract metadata from a YouTube video

        Args:
            youtube_id (str): YouTube video ID

        Returns:
            dict: Dictionary containing video metadata
        """
        logger.info(f"Extracting metadata for YouTube ID: {youtube_id}")

        youtube_url = f"https://www.youtube.com/watch?v={youtube_id}"

        ydl_opts = {
            'format': 'bestaudio/best',
            'writeinfojson': True,
            'writesubtitles': False,
            'writeannotations': False,
            'noplaylist': True,
            'quiet': True,
            'no_warnings': True,
            'ignoreerrors': False,
            'outtmpl': f'{self.temp_dir}/{youtube_id}/%(id)s.%(ext)s',
        }

        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            info = ydl.extract_info(youtube_url, download=False)

            # Map the yt-dlp info to our model structure
            metadata = {
                'title': info.get('title', ''),
                'full_title': info.get('fulltitle', ''),
                'description': info.get('description', ''),
                'youtube_id': youtube_id,
                'duration': info.get('duration', 0),
                'duration_string': info.get('duration_string', "00:00"),
                'thumbnail': info.get('thumbnail', ''),
                'tags': info.get('tags', []),
                'categories': info.get('categories', []),
            }

        logger.info(f"Metadata extraction completed for {youtube_id}")
        return metadata

    def extract_transcript(self, youtube_id):
        """
        Extract transcript from a YouTube video

        Args:
            youtube_id (str): YouTube video ID

        Returns:
            tuple: (language_code, list of transcript segments)
        """
        logger.info(f"Extracting transcript for YouTube ID: {youtube_id}")

        youtube_url = f"https://www.youtube.com/watch?v={youtube_id}"

        # Create a temporary directory for this transcription
        with tempfile.TemporaryDirectory(prefix=f"shadowify_transcript_{youtube_id}_") as temp_dir:
            logger.info(
                f"Created temporary directory for transcript: {temp_dir}")

            audio_file = f"{temp_dir}/{youtube_id}.wav"

            # Download audio from YouTube
            ydl_opts = {
                'format': 'm4a/bestaudio/best',
                'postprocessors': [{  # Extract audio using ffmpeg
                    'key': 'FFmpegExtractAudio',
                    'preferredcodec': 'wav',
                }],
                'outtmpl': f'{temp_dir}/{youtube_id}.%(ext)s',
                'quiet': True,
                'no_warnings': True,
                'noplaylist': True,
            }

            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                logger.info(f"Downloading audio for {youtube_id}")
                ydl.download([youtube_url])

            # Transcribe audio using Whisper
            logger.info(f"Transcribing audio for {youtube_id}")
            segments, info = self.whisper_model.transcribe(
                audio_file,
                beam_size=5,
                word_timestamps=True,
                vad_filter=True
            )

            language_code = info.language
            transcript_segments = []

            for segment in segments:
                transcript_segments.append({
                    'start': segment.start,
                    'end': segment.end,
                    'text': segment.text.strip()
                })

            logger.info(f"Transcription completed for {youtube_id}")
            # Temporary directory will be automatically cleaned up after this block

        return language_code, transcript_segments

    def _format_duration(self, seconds):
        """
        Format duration in seconds to a human-readable string (HH:MM:SS)

        Args:
            seconds (int): Duration in seconds

        Returns:
            str: Formatted duration string
        """
        if not seconds:
            return "00:00"

        hours = seconds // 3600
        minutes = (seconds % 3600) // 60
        seconds = seconds % 60

        if hours > 0:
            return f"{hours:02d}:{minutes:02d}:{seconds:02d}"
        else:
            return f"{minutes:02d}:{seconds:02d}"

    def _get_best_thumbnail(self, info):
        """
        Get the best quality thumbnail URL from video info

        Args:
            info (dict): Video info from yt-dlp

        Returns:
            str: URL of the best thumbnail
        """
        if not info or 'thumbnails' not in info or not info['thumbnails']:
            return ""

        # Sort thumbnails by resolution (width * height) in descending order
        sorted_thumbnails = sorted(
            info['thumbnails'],
            key=lambda x: (x.get('width', 0) or 0) * (x.get('height', 0) or 0),
            reverse=True
        )

        # Return the URL of the highest resolution thumbnail
        return sorted_thumbnails[0].get('url', "")
