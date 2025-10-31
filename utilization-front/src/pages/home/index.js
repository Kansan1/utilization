import React, { Component } from "react";
import "./home.scss";
import DeviceForm from "./看板管理/DeviceForm";
import Dashboard from "./设备看板";
import EquipmentMaintenance from "./设备维护维修看板";
import EquipmentSchedule from "./设备维护维修看板管理/EquipmentSchedule";
import MyTabs from "./设备维护维修看板管理";
// import EquipmentMaintenance from "./设备维护维修看板/ScrollTable";
// 导入StorageScreen组件

class Home extends Component {
  constructor(props) {
    super(props);
    this.headerRef = React.createRef();
    this.hideHeaderTimer = null;
  }

  renderKey = 1;

  state = {
    // 菜单数据
    menuItems: [
      { id: 1, title: "设备管理看板", icon: "home" },
      { id: 2, title: "设备管理后台", icon: "order" },
       { id: 3, title: "设备维护维修看板", icon: "home" },
      { id: 4, title: "设备维护维修看板管理", icon: "order" },
      /*{ id: 4, title: "4库", icon: "product" },
      { id: 5, title: "5库", icon: "user" },
      { id: 6, title: "6库", icon: "setting" },
      { id: 7, title: "7库", icon: "setting" }, */
    ],
    activeMenu: 3, // 当前选中的菜单项
    contentTitle: "设备看板", // 右侧内容区标题
    isFullscreen: false, // 是否全屏显示
    showHeader: true, // 控制顶部栏显示/隐藏
    currentTime: new Date(),
  };

  // 在组件挂载后添加事件监听
  componentDidMount() {
    // 每秒更新一次时间
    this.timer = setInterval(() => {
      this.setState({
        currentTime: new Date(),
      });
    }, 1000);

    setInterval(() => {
      this.renderKey *= -1;
    }, 3600000);

    //21600000
  }

  // 组件卸载时清理
  componentWillUnmount() {
    if (this.hideHeaderTimer) {
      clearTimeout(this.hideHeaderTimer);
    }
    if (this.timer) {
      clearInterval(this.timer); // 清除时间更新定时器
    }
  }

  // 处理鼠标进入顶部区域
  handleMouseEnterTopArea = () => {
    // 取消任何待执行的隐藏计时器
    if (this.hideHeaderTimer) {
      clearTimeout(this.hideHeaderTimer);
      this.hideHeaderTimer = null;
    }

    if (this.state.isFullscreen) {
      this.setState({ showHeader: true });
    }
  };

  // 处理鼠标离开顶部区域
  handleMouseLeaveTopArea = () => {
    // 只在全屏模式下处理
    if (this.state.isFullscreen) {
      // 设置一个延迟，让用户有时间移动到顶部栏
      this.hideHeaderTimer = setTimeout(() => {
        this.setState({ showHeader: false });
      }, 300); // 300ms延迟，可以根据需要调整
    }
  };

  // 处理菜单点击
  handleMenuClick = (menuId, title) => {
    this.setState({
      activeMenu: menuId,
      contentTitle: title,
    });
    // 这里可以根据菜单ID加载不同内容
  };

  // 切换全屏模式
  toggleFullscreen = () => {
    this.setState((prevState) => ({
      isFullscreen: !prevState.isFullscreen,
      showHeader: true, // 切换时总是先显示头部
    }));

    // 如果进入全屏模式，设置延时自动隐藏顶部栏
    if (!this.state.isFullscreen) {
      setTimeout(() => {
        if (this.state.isFullscreen) {

          this.setState({ showHeader: false });
        }
      }, 2000); // 2秒后自动隐藏
    }
  };

  // 格式化时间
  formatTime = (date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    const hours = String(date.getHours()).padStart(2, "0");
    const minutes = String(date.getMinutes()).padStart(2, "0");
    const seconds = String(date.getSeconds()).padStart(2, "0");

    // 获取星期几
    const weekdays = [
      "星期日",
      "星期一",
      "星期二",
      "星期三",
      "星期四",
      "星期五",
      "星期六",
    ];
    const weekDay = weekdays[date.getDay()];

    return `${year}-${month}-${day} ${weekDay} ${hours}:${minutes}:${seconds}`;
  };


  // 渲染当前选中菜单对应的内容
  renderContent() {
    const { activeMenu  } = this.state;

    switch (activeMenu) {
      case 1:
        return (
            <Dashboard
                isFullscreen={this.isFullscreen}
                currentTime={this.formatTime(this.state.currentTime)}
                renderKey={this.renderKey}
            />
        );
        case 2:
          return (
                <DeviceForm  />
          );
      case 3:
        return (
            <EquipmentMaintenance
                isFullscreen={this.isFullscreen}
                currentTime={this.formatTime(this.state.currentTime)}
                  />
        );
      case 4:
        return (
            <MyTabs
            />
        );
      // case 2:
      //   return (
      //     <div className="default-content">
      //       <h3>2库内容</h3>
      //       <p>这里是2库的内容展示区域</p>
      //     </div>
      //   );
      // case 3:
      //   return (
      //     <div className="default-content">
      //       <h3>3库内容</h3>
      //       <p>这里是3库的内容展示区域</p>
      //     </div>
      //   );
      // case 4:
      //   return (
      //     <div className="default-content">
      //       <h3>4库内容</h3>
      //       <p>这里是4库的内容展示区域</p>
      //     </div>
      //   );
      // case 5:
      //   return (
      //     <div className="default-content">
      //       <h3>5库内容</h3>
      //       <p>这里是5库的内容展示区域</p>
      //     </div>
      //   );
      // case 6:
      //   return (
      //     <div className="default-content">
      //       <h3>6库内容</h3>
      //       <p>这里是6库的内容展示区域</p>
      //     </div>
      //   );
      // case 7:
      //   return (
      //     <div className="default-content">
      //       <h3>7库内容</h3>
      //       <p>这里是7库的内容展示区域</p>
      //     </div>
      //   );
      default:
        return (
          <div className="welcome-message">
            <h3>欢迎使用仓库管理系统</h3>
            <p>请从左侧选择功能菜单</p>
          </div>
        );
    }
  }

  render() {
    const {
      menuItems,
      activeMenu,
      contentTitle,
      isFullscreen,
      showHeader,
      currentTime,
    } = this.state;

    return (
      <div className="home-container" >
        {!isFullscreen && (
          <div className="sidebar">
            <div className="logo">
              <h2>工业技术大屏系统</h2>
            </div>
            <ul className="menu-list">
              {menuItems.map((item) => (
                <li
                  key={item.id}
                  className={activeMenu === item.id ? "active" : ""}
                  onClick={() => this.handleMenuClick(item.id, item.title)}
                >
                  <span className={`icon icon-${item.icon}`}></span>
                  <span className="menu-title">{item.title}</span>
                </li>
              ))}
            </ul>
          </div>
        )}

        <div   className={`content-area ${isFullscreen ? "fullscreen" : ""}`}>
          {isFullscreen && (
            <div
              className="header-sensor"
              onMouseEnter={this.handleMouseEnterTopArea}
              onMouseLeave={this.handleMouseLeaveTopArea}
            ></div>
          )}

          <div
            className={`content-header ${
              isFullscreen && !showHeader ? "hidden" : ""
            }`}
            ref={this.headerRef}
            onMouseEnter={this.handleMouseEnterTopArea}
            onMouseLeave={this.handleMouseLeaveTopArea}

          >
            <h1 className="content-title">{contentTitle}</h1>
            <div className="header-actions">
              <button
                className="fullscreen-btn"
                onClick={this.toggleFullscreen}
                title={isFullscreen ? "退出全屏" : "全屏显示"}
              >
                {isFullscreen ? "退出全屏" : "全屏显示"}
              </button>
            </div>
          </div>
          {this.renderContent()}

        </div>
      </div>
    );
  }
}

export default Home;
