import React, { Component } from "react";
import "./mouseTrail.scss";

class MouseTrail extends Component {
  particles = [];

  componentDidMount() {
    // 添加鼠标移动事件监听
    document.addEventListener("mousemove", this.handleMouseMove);

    // 定期清理粒子
    this.cleanupInterval = setInterval(() => {
      this.particles.forEach((particle) => {
        if (particle.classList.contains("fade") && particle.parentNode) {
          particle.remove();
          this.particles.splice(this.particles.indexOf(particle), 1);
        }
      });
    }, 1000);
  }

  componentWillUnmount() {
    // 清除事件监听和定时器
    document.removeEventListener("mousemove", this.handleMouseMove);
    clearInterval(this.cleanupInterval);
  }

  handleMouseMove = (e) => {
    // 在鼠标位置创建粒子
    this.createParticle(e.clientX, e.clientY);
  };

  createParticle = (x, y) => {
    const trailContainer = document.getElementById("trail-container");
    if (!trailContainer) return;

    const particle = document.createElement("div");
    particle.classList.add("particle");
    particle.style.left = `${x}px`;
    particle.style.top = `${y}px`;
    particle.style.backgroundColor = this.getRandomColor();
    trailContainer.appendChild(particle);
    this.particles.push(particle);

    // 添加淡出动画
    setTimeout(() => {
      particle.classList.add("fade");
      particle.addEventListener("animationend", () => {
        if (particle.parentNode) {
          particle.remove();
          this.particles.splice(this.particles.indexOf(particle), 1);
        }
      });
    }, 100);
  };

  getRandomColor = () => {
    const letters = "0123456789ABCDEF";
    let color = "#";
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  };

  render() {
    return <div id="trail-container"></div>;
  }
}

export default MouseTrail;
