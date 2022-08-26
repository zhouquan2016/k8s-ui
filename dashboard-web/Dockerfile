FROM node as build
WORKDIR /app
ADD . .
# RUN npm config set registry https://npmmirror.com/mirrors/
RUN npm install 
RUN npm run build

FROM nginx as deploy
COPY --from=build /app/build/** /usr/share/nginx/html/