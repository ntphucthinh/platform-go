FROM node:23-alpine

RUN apk update
# openapi-generator-cli
RUN apk add --no-cache openjdk21-jre-headless

RUN npm install -g npm
# openapi cli
RUN npm install -g @openapitools/openapi-generator-cli
# 
RUN openapi-generator-cli version

# Add openapi2aspida for make codegen
RUN npm install -g openapi2aspida typescript ts-node

WORKDIR /app

# Keep container running in background
CMD ["sh"]
