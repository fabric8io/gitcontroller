FROM centos:7

RUN yum install -y git

ENV PATH $PATH:/usr/local/gitcontroller/

ADD ./build/gitcontroller /usr/local/gitcontroller/

CMD gitcontroller run
