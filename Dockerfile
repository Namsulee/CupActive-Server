FROM scratch
ADD cupactive-server /
ADD ./www-static/ /www-static/
CMD ["/cupactive-server"]
