sudo docker run --rm \
  -u "$(id -u):$(id -g)" \
  -v ./books:/data \
  ebook-pdf \
  /data -o /data -r