FROM debian:stable-slim
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y
RUN apt-get install apt-transport-https -y
RUN apt-get install apt-utils -y
RUN apt-get install gcc -y
RUN apt-get install g++ -y
RUN apt-get install nano -y
RUN apt-get install tar -y
RUN apt-get install bash -y
RUN apt-get install sudo -y
RUN apt-get install openssl -y
RUN apt-get install git -y
RUN apt-get install make -y
RUN apt-get install wget -y
RUN apt-get install curl -y
RUN apt-get install net-tools -y
RUN apt-get install iproute2 -y
RUN apt-get install bc -y
RUN apt-get install libasound2-dev -y
RUN apt-get install libhidapi-dev -y
#RUN apt-get install udev -y
RUN apt-get install pkg-config -y
RUN apt-get install alsa-utils -y

#RUN amixer sset Master 100%
#RUN amixer sset Master unmute

# Setup StreamDeck
#RUN mkdir -p /etc/udev/rules.d
#RUN echo 'SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006d", MODE="0666", GROUP="morphs"' > /etc/udev/rules.d/99-streamdeck.rules
#RUN udevadm control --reload-rules
#RUN udevadm trigger

# Setup User
ENV TZ="US/Eastern"
ARG USERNAME="morphs"
ARG PASSWORD="asdf"
RUN useradd -m $USERNAME -p $PASSWORD -s "/bin/bash"
RUN mkdir -p /home/$USERNAME
RUN chown -R $USERNAME:$USERNAME /home/$USERNAME
RUN usermod -aG sudo $USERNAME
RUN echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
RUN usermod -a -G dialout $USERNAME
RUN usermod -a -G audio $USERNAME
#RUN chown morphs:morphs /dev/streamdeck
USER $USERNAME

# install go with specific version and progress
WORKDIR /home/$USERNAME
COPY ./go_install.sh /home/$USERNAME/go_install.sh
RUN sudo chmod +x /home/$USERNAME/go_install.sh
RUN sudo chown $USERNAME:$USERNAME /home/$USERNAME/go_install.sh
RUN /home/$USERNAME/go_install.sh
RUN sudo tar --checkpoint=100 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo ZXhwb3J0IEdPX1RBUl9LSUxPQllURVM9JChwcmludGYgIiUuM2ZcbiIgJChlY2hvICIkKHN0YXQgLS1mb3JtYXQ9IiVzIiAvaG9tZS9tb3JwaHMvZ28udGFyLmd6KSAvIDEwMDAiIHwgYmMgLWwpKSAmJiBlY2hvIEV4dHJhY3RpbmcgWyRUQVJfQ0hFQ0tQT0lOVF0gb2YgJEdPX1RBUl9LSUxPQllURVMga2lsb2J5dGVzIC91c3IvbG9jYWwvZ28= | base64 -d ; echo); eval $cmd"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
RUN echo "PATH=$PATH:/usr/local/go/bin" | tee -a /home/$USERNAME/.bashrc

# prep files for building server binary
USER root
RUN mkdir /home/morphs/StreamDeckServer
COPY . /home/morphs/StreamDeckServer
RUN chown -R $USERNAME:$USERNAME /home/morphs/StreamDeckServer
USER $USERNAME
WORKDIR /home/$USERNAME

# build server binary
ARG GO_ARCH=amd64
WORKDIR StreamDeckServer
RUN /usr/local/go/bin/go mod tidy
RUN GOOS=linux GOARCH=$GO_ARCH /usr/local/go/bin/go build -o /home/morphs/StreamDeckServer/server
#ENTRYPOINT [ "/home/morphs/StreamDeckServer/server" ]
#ENTRYPOINT [ "bash" ]

USER root
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER $USERNAME
ENTRYPOINT ["/entrypoint.sh"]
