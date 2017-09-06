for APP in ingress frontend order user item; do
    cat <<EOF > ${APP}-csr.json 
{
  "CN": "${APP}",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "US",
      "L": "Atlanta",
      "O": "apps",
      "OU": "example",
      "ST": "Georgia"
    }
  ]
}
EOF

    cfssl gencert \
        -ca=ca.pem \
        -ca-key=ca-key.pem \
        -config=ca-config.json \
        -hostname=${APP},${APP}.default,${APP}.svc.cluster.local \
        -profile=kubernetes \
        ${APP}-csr.json | cfssljson -bare ${APP}
done
