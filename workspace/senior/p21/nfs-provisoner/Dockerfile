FROM redhat/ubi8-minimal:latest

COPY ./nfs-provisioner /nfs-provisioner

RUN chmod +x /nfs-provisioner

ENTRYPOINT ["/nfs-provisioner"]

