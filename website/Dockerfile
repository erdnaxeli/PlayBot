FROM php:5-apache
RUN ln -s /etc/apache2/mods-available/rewrite.load /etc/apache2/mods-enabled/rewrite.load && \
    docker-php-ext-install pdo_mysql
#COPY ./config.php /var/www/html/
COPY app/ statics/ .htaccess /var/www/html/
