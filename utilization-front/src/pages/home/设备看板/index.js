import React from 'react';
import UtilizationRateDay from './utilization_rate_day';
import UtilizationRateMonth from './utilization_rate_month';
import InspectionDay from './设备点检';
import MaintenanceMonth from './设备维护完成率';
import RepairMonth from './设备维修次数';
import {homeAPi as api, homeAPi} from "../../../api";
import {message} from "antd";
import {center} from "@antv/g2plot/lib/plots/sankey/sankey";

const Dashboard = ({ isFullscreen, currentTime, renderKey }) => {
    return (
        <div className={`content-body ${isFullscreen ? "fullscreen" : ""} ` } >
            <div className="header-container">
                <div className="left-placeholder"></div>
                <div className="content-title"  >设备管理看板</div>
                <div className="time-display">{currentTime}</div>
            </div>
            <div style={{ marginLeft: '24px', marginRight: '24px' }} key={renderKey}>
                <div className="library-one-content">
                    <div className="top-row">
                        <div className="chart-section">
                            <UtilizationRateDay />
                        </div>
                        <div className="chart-section">
                            <UtilizationRateMonth />
                        </div>
                        <div className="chart-section">
                            <InspectionDay />
                        </div>
                    </div>
                    <div className="bottom-row">
                        <div className="chart-section">
                            <MaintenanceMonth />
                        </div>
                        <div className="chart-section">
                            <RepairMonth />
                        </div>

                    </div>
                </div>
            </div>
        </div>
    );
};



export default Dashboard;
