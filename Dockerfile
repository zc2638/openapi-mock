FROM centos
WORKDIR /app
ENV PATH /app:$PATH
COPY mock /app/mock
RUN chmod +x mock
CMD ["mock"]