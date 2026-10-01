import React from 'react';
import { Chart as ChartJS, ArcElement, Tooltip } from 'chart.js';
import { Doughnut } from 'react-chartjs-2';


const options={
    legend: {
        display: false,
    },
};


ChartJS.register(ArcElement, Tooltip);

const textCenterPlugin = {
    id: 'textCenter',
    beforeDraw(chart) {
        const { ctx, width, height } = chart;
        ctx.restore();

        // Configuration de la police (s'adapte dynamiquement à la taille du graphique)
        const fontSize = (height / 120).toFixed(2);
        ctx.font = `${fontSize}em sans-serif`;
        ctx.textBaseline = 'middle';
        ctx.fillStyle = '#1e293b'; // Couleur du texte

        var sum = 0;
        for (var i = 0; i < chart.config.data.datasets[0].data.length; i++) {
            sum += chart.config.data.datasets[0].data[i];
        }
        // Texte à afficher
        const text = sum;

        // Calcul du centrage horizontal et vertical
        const textX = Math.round((width - ctx.measureText(text).width) / 2);
        const textY = height / 2;

        // Dessiner le texte
        ctx.fillText(text, textX, textY);
        ctx.save();
    },
};



class DonutWithText extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: {labels:[],datasets:[]}
        };
    }

    static fixData (alertes) {
        var ndata = {labels:[],datasets:[]};
        
        var o = {data:[]};
        for (var i in alertes){
            var ni = alertes[i];
            for (var k in ni){
                //if (ni.hasOwnProperty(k)) {
                if (Object.prototype.hasOwnProperty.call(ni,k)) {
                    ndata.labels.push(k);
                    o.data.push(ni[k]);
                }
            }
        }

        o.backgroundColor = [ '#ff7f0e', '#1f77b4', '#aec7e8', '#ffca28', '#d4e157','#4caf50','#26a69a','#00e5ff', '#00b0ff', '#ff1744' ]; 
        o.hoverBackgroundColor = [ '#ff9f2e', '#3f97d4', '#bed7f8', '#ffca28', '#d4e157','#4caf50','#26a69a','#00e5ff', '#00b0ff', '#ff1744'];
        ndata.datasets.push(o);

        return ndata;
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        if (nextProps.data !== null && nextProps.data.length !== 0 ) {
            return {
                data: DonutWithText.fixData(nextProps.data)
            };
        }
        return null;
    }


    render() {
        var a = this.props.data;
        if (a === null || a.length === 0) {
            return ( 
                <div>...</div> 
            );
        }
        return (
            <div>
                {this.props.title} : 
                <Doughnut data={this.state.data} options={options} plugins={[textCenterPlugin]}  height={150} width={180} />
            </div>
        );
    }
}

export default DonutWithText;
