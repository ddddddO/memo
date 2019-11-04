FROM ruby:2.6-alpine

WORKDIR /web

COPY /web .
COPY /deployments/dockerfiles/web/hosts .

EXPOSE 8080

RUN cat hosts >> /etc/hosts \
 && rm Gemfile.lock \
 && rm -rf vendor \
 && sed -i -e 's/source "https:\/\/rubygems.org"/source "http:\/\/rubygems.org"/g' ./Gemfile \
 && bundle install --path vendor/bundle

ENTRYPOINT ["bundle exec unicorn"]
#CMD ["-c config/unicorn.rb"]
