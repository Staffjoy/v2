set -e

sudo apt install -y -q  nginx

echo '
server {
    listen 80;
    server_name kubernetes.staffjoy-v2.local;
    location / {
        proxy_pass http://localhost:8001;
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

' | sudo tee /etc/nginx/sites-enabled/default

sudo service nginx restart
