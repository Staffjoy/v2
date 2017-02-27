set -e

apt-get -y --force-yes install nginx
echo '
server {
    listen 80; 
    server_name kubernetes.staffjoy-v2.local;
    location / {
        proxy_pass http://localhost:8080;
    }
}
server {
    listen 80; 
    server_name *.staffjoy-v2.local staffjoy-v2.local;
    location / { 
        proxy_pass http://10.0.0.99:80; 
        proxy_set_header Host            $host; 
        proxy_set_header X-Forwarded-For $remote_addr; 
    }
} 

'  > /etc/nginx/sites-enabled/default
service nginx restart
