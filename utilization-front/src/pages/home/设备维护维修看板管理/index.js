import React, { useState } from 'react';
import { Tabs } from 'antd';


import EquipmentSchedule from "./EquipmentSchedule";
import EquipmentScheduleWithCategoryType from "./EquipmentScheduleType";
import EquipmentPlanSpecific from "./EquipmentPlanSpecific";
import DailyRepairTask from "./DailyRepairTask";

const { TabPane } = Tabs;

export default function MyTabs() {
    const [activeKey, setActiveKey] = useState('1');

    return (
        <Tabs  style={{ marginTop: '10px', marginLeft: '10px', minHeight: '500px', backgroundColor: 'transparent' }} activeKey={activeKey} onChange={setActiveKey}>

            <TabPane tab="日维修任务" key="1">
                <DailyRepairTask />
            </TabPane>
            <TabPane tab="具体维护计划" key="2">
                <EquipmentPlanSpecific />
            </TabPane>
            <TabPane tab="维护计划" key="3">
                <EquipmentSchedule />
            </TabPane>
            <TabPane tab="设备类型" key="4">
                <EquipmentScheduleWithCategoryType />
            </TabPane>
        </Tabs>
    );
}