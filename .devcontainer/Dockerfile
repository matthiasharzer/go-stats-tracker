FROM mcr.microsoft.com/devcontainers/go:2-1.25-trixie

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        libtesseract-dev \
        libleptonica-dev \
        tesseract-ocr-eng \
    		tesseract-ocr-deu \
        tesseract-ocr \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
