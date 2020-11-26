docker run -it --rm \
      -v /opt/app/certs:/etc/letsencrypt \
      -v /opt/app/certs-data:/data/letsencrypt \
      deliverous/certbot \
      certonly \
      --webroot --webroot-path=/data/letsencrypt \
      -d api.balconygames.com

docker run -it --rm \
      -v /opt/app/certs:/etc/letsencrypt \
      -v /opt/app/certs-data:/data/letsencrypt \
      deliverous/certbot \
      certonly \
      --webroot --webroot-path=/data/letsencrypt \
      -d metabase.balconygames.com

docker run -it --rm \
      -v /opt/app/certs:/etc/letsencrypt \
      -v /opt/app/certs-data:/data/letsencrypt \
      deliverous/certbot \
      certonly \
      --webroot --webroot-path=/data/letsencrypt \
      -d dashboard.balconygames.com
