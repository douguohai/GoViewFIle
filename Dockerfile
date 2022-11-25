FROM centos:centos7
COPY fonts/* /usr/share/fonts/ChineseFonts/


# 设置固定的项目路径
ENV WORKDIR /var/www/GoViewFile
ENV TZ=Asia/Shanghai

RUN yum update -y  && yum  install -y deltarpm  wget   && export DISPLAY=:0.0 &&\
    yum install -y libreoffice  libreoffice-headless  libreoffice-writer ImageMagick &&\
    wget https://douguohai.oss-cn-shenzhen.aliyuncs.com/shard/wkhtmltox-0.12.6.1-2.almalinux8.aarch64.rpm  &&\
    rpm --rebuilddb && yum install -y openssl && yum install -y xorg-x11-fonts-75dpi &&\
    rpm -ivh wkhtmltox-0.12.6.1-2.almalinux8.aarch64.rpm && localedef -i en_US -f UTF-8 en_US.UTF-8 && \


# 添加应用可执行文件，并设置执行权限
ADD main   $WORKDIR/main
RUN chmod +x $WORKDIR/main

# 添加I18N多语言文件、静态文件、配置文件、模板文件
ADD public   $WORKDIR/public
ADD config   $WORKDIR/config
ADD template $WORKDIR/template

# 添加本地上传文件目录
COPY cache/convert/  $WORKDIR/cache/convert/
COPY cache/download/  $WORKDIR/cache/download/
COPY cache/local/  $WORKDIR/cache/local/
COPY cache/pdf/  $WORKDIR/cache/pdf/
# jar包，用于将.msg文件转eml文件
COPY library/emailconverter-2.5.3-all.jar   /usr/local/emailconverter-2.5.3-all.jar
#pdf 添加水印
COPY library/pdfcpu    /usr/local/bin/pdfcpu

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
# 如果需要进入容器调式，可以注释掉下面的CMD. 
CMD  ./main  


# ------------------------------------本地打包镜像---------------------
# docker build -t  goviewfile:v0.7  .
# docker run -d  -p 8082:8082 镜像ID
