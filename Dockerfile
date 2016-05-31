FROM centos:7

ENV PATH $PATH:/usr/local/gitcontroller/

ADD ./build/gitcontroller /usr/local/gitcontroller/

CMD gitcontroller run
