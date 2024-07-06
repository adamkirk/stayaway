FROM php:8.2-cli AS base

RUN apt update

# For composer
RUN apt install -y zip unzip

RUN pecl install openswoole
RUN docker-php-ext-enable openswoole

RUN docker-php-ext-configure pcntl --enable-pcntl && docker-php-ext-install pcntl

RUN pecl install redis && docker-php-ext-enable redis
RUN apt install -y libpq-dev && docker-php-ext-install pdo pdo_pgsql

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

FROM base as dev

RUN apt install -y nodejs npm

FROM base as prod

RUN apt clean