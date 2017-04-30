FROM node:7.9.0-alpine
WORKDIR .
ADD ./frontend .
RUN npm install
EXPOSE 1337
CMD ["npm", "run", "start:prod"]
