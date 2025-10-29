openssl rand -base64 48 | tr -dc 'A-Za-z0-9' | fold -w 32 | head -n 1 > rabbitmq_pass.txt
chmod 600 rabbitmq_pass.txt
touch .env
echo -e "RABBITMQ_DEFAULT_USER=admin" >> .env
echo -e "RABBITMQ_DEFAULT_PASS=$(cat rabbitmq_pass.txt)" >> .env