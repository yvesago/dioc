import React from 'react';
import {Doughnut, Chart} from 'react-chartjs-2';

// some of this code is a variation on https://jsfiddle.net/cmyker/u6rr5moq/
var originalDoughnutDraw = Chart.controllers.doughnut.prototype.draw;
Chart.helpers.extend(Chart.controllers.doughnut.prototype, {
    draw: function() {
        originalDoughnutDraw.apply(this, arguments);
    
        var chart = this.chart;
        var width = chart.chart.width,
            height = chart.chart.height,
            ctx = chart.chart.ctx;

        var fontSize = (height / 114).toFixed(2);
        ctx.font = fontSize + 'em sans-serif';
        ctx.textBaseline = 'middle';

        var sum = 0;
        for (var i = 0; i < chart.config.data.datasets[0].data.length; i++) {
            sum += chart.config.data.datasets[0].data[i];
        }

        var text = sum,
            textX = Math.round((width - ctx.measureText(text).width) / 2),
            textY = height / 2;

        ctx.fillText(text, textX, textY);
    }
});

const options={
    legend: {
        display: false,
    },
};

class DonutWithText extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: {labels:[],datasets:[]}
        };
    }

    componentWillReceiveProps(props) {
        if (props.data.length !== 0) {
            this.setState({ alertes: props.data });
            this.setData(props.data);
        }
    }

    setData(alertes) {
        var ndata = {labels:[],datasets:[]};
        
        var o = {data:[]};
        for (var i in alertes){
            var ni = alertes[i];
            for (var k in ni){
                if (ni.hasOwnProperty(k)) {
                    ndata.labels.push(k);
                    o.data.push(ni[k]);
                }
            }
        }

        o.backgroundColor = [ '#ff7f0e', '#1f77b4', '#aec7e8', '#ffca28', '#d4e157','#4caf50','#26a69a','#00e5ff', '#00b0ff', '#ff1744' ]; 
        o.hoverBackgroundColor = [ '#ff4f00', '#3f97d4', '#bed7f8', '#ffca28', '#d4e157','#4caf50','#26a69a','#00e5ff', '#00b0ff', '#ff1744'];
        ndata.datasets.push(o);

        this.setState({ data: ndata });
    }



    render() {
        var a = this.props.data;
        if (a.length === 0) {
            return ( 
                <div>...</div> 
            );
        }
        return (
            <div>
                {this.props.title} : 
                <Doughnut data={this.state.data} options={options}  height={150} width={180} />
            </div>
        );
    }
}

export default DonutWithText;
