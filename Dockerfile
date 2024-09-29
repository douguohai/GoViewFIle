FROM golang:alpine AS builder

WORKDIR /root
ENV WORKDIR=/root
ADD . $WORKDIR
RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w GO111MODULE=on  && go install github.com/pdfcpu/pdfcpu/cmd/pdfcpu@latest  && go build -o GoViewFIle


FROM fedora:latest AS runner
# 设置固定的项目路径
ENV TZ=Asia/Shanghai
# 设置固定的项目路径
WORKDIR /var/www/GoViewFile
ENV WORKDIR=/var/www/GoViewFile
ENV DISPLAY=:0.0

RUN yum makecache  &&\
    yum install -y deltarpm  wget  libreoffice  libreoffice-headless  libreoffice-writer ImageMagick  openssl  xorg-x11-fonts-75dpi && yum clean all

COPY --from=builder /root/fonts/* /usr/share/fonts/ChineseFonts/

# 添加I18N多语言文件、静态文件、配置文件、模板文件å
COPY --from=builder /root/public/  $WORKDIR/public
COPY --from=builder /root/config/  $WORKDIR/config
COPY --from=builder /root/template/ $WORKDIR/template
COPY --from=builder /go/bin/pdfcpu /usr/local/bin/pdfcpu
# 添加应用可执行文件，并设置执行权限
COPY --from=builder /root/GoViewFIle   $WORKDIR/go_view_file

RUN  cd $WORKDIR && chmod +x $WORKDIR/go_view_file && mkdir cache && cd cache && mkdir convert download local pdf

# 如果需要进入容器调式，可以注释掉下面的CMD.
CMD  ["./go_view_file"]


# ------------------------------------本地打包镜像---------------------
# docker build -t  goviewfile:v0.7  .
# docker run -d  -p 8082:8082 镜像ID
