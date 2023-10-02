
<template>
    <div>
        <div id="main" style="width: 1200px; height: 760px"></div>
    </div>
</template>
<script setup>
import axios from 'axios';
import * as echarts from 'echarts';
import { nextTick } from "vue";
axios.defaults.withCredentials = true


const service = axios.create({
  baseURL: "/api",
  timeout: 99999
})
let acitveAxios = 0
let timer

// http request 拦截器
service.interceptors.request.use(
  config => {

    config.headers = {
      'Content-Type': 'application/json',
      ...config.headers
    }
    return config
  }
)


const test = async () => {
    await nextTick();
    var data = {}
    await service.get("/api/visitor/count").then(function (response) {
        data = response.data.data;
        console.log(response);
    }).catch(function (error) {
        console.log(error);
    });

    console.log(data.data);
    var chartDom = document.getElementById('main');
    var myChart = echarts.init(chartDom);
    var option;

    option = {
        title: {
            text: '每天在线统计 阈值低于2标红',
            subtext: '人数'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'cross'
            }
        },
        toolbox: {
            show: true,
            feature: {
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            // prettier-ignore
            data: data.x
        },
        yAxis: {
            type: 'value',
            axisLabel: {
                formatter: '{value} 人'
            },
            axisPointer: {
                snap: true
            }
        },
        visualMap: {
            show: false,
            dimension: 0,
            pieces: data.pieces
        },
        series: [
            {
                name: 'Electricity',
                type: 'line',
                smooth: true,
                // prettier-ignore
                data: data.data,
                markArea: {
                    itemStyle: {
                        color: 'rgba(255, 173, 177, 0.4)'
                    },
                    data: [
                        // [
                        //     {
                        //         name: 'Morning Peak',
                        //         xAxis: '07:30'
                        //     },
                        //     {
                        //         xAxis: '10:00'
                        //     }
                        // ],
                        // [
                        //     {
                        //         name: 'Evening Peak',
                        //         xAxis: '17:30'
                        //     },
                        //     {
                        //         xAxis: '21:15'
                        //     }
                        // ]
                    ]
                }
            }
        ]
    };

    option && myChart.setOption(option);

}
test()

</script>
