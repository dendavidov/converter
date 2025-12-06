FROM python:3.12-slim

ENV DEBIAN_FRONTEND=noninteractive \
    QTWEBENGINE_DISABLE_SANDBOX=1 \
    QTWEBENGINE_CHROMIUM_FLAGS="--disable-gpu --no-sandbox"

RUN apt-get update && \
    apt-get install -y --no-install-recommends calibre dbus-x11 fontconfig && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY ebook_to_pdf.py /app/ebook_to_pdf.py
RUN chmod 755 /app/ebook_to_pdf.py

ENTRYPOINT ["python", "/app/ebook_to_pdf.py"]