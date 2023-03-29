FROM scratch
COPY wfind /wfind
ENTRYPOINT ["/wfind"]

