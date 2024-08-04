FROM registry.cmd.navi-tech.in/common/node:18.15.0-alpine3.16 as build

WORKDIR /app
COPY . /app
RUN yarn install
RUN yarn build

FROM registry.cmd.navi-tech.in/common/nginx:1.23.3-alpine
COPY --from=build /app/dist /usr/share/nginx/html
RUN adduser -u 4000 non-root-user -D ''
RUN chown -R 4000:4000 /var/cache/nginx \
    && chown -R 4000:4000 /etc/nginx/conf.d/ \
    && chown -R 4000:4000 /usr/share/nginx \
    && chmod -R g+w /var/cache/nginx \
    && touch /var/run/nginx.pid \
    && chown -R 4000:4000 /var/run/nginx.pid \
    && ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log
RUN rm /etc/nginx/conf.d/default.conf

USER 4000

COPY nginx/nginx.conf /etc/nginx/conf.d
COPY entrypoint.sh /

USER 0

RUN chmod +x entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]

USER 4000

CMD ["nginx", "-g", "daemon off;"]
