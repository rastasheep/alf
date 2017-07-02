FROM node:8-alpine

MAINTAINER Aleksandar Diklic "https://github.com/rastasheep"

RUN yarn global add @angular/cli@1.2.0 \
  && ng set --global packageManager=yarn \

CMD ["ng", "serve"]
