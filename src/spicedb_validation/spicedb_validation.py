import os
import sys

import ctypes

dll = ctypes.cdll.LoadLibrary(
    os.path.join(os.path.dirname(__file__), "dll", "spicedb_validation.so")
)


def validate_url(url: str):
    dll.validateURL(url.encode("utf-8"))
