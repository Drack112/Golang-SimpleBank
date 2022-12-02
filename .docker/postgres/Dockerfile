FROM postgres:latest as development

RUN usermod -u 1000 postgres

CMD ["-p", "5433"]

FROM postgres:latest as test

RUN usermod -u 1000 postgres
