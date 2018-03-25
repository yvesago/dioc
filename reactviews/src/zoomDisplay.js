import React from 'react';
import ReactDOM from 'react-dom';
import { MapControl } from 'react-leaflet';
import L from 'leaflet';

import './react.leaflet.posdisplay.css';


export default class ZoomDisplay extends MapControl {
    componentWillMount() {
        const legend = L.control({position: 'bottomleft'});
        const jsx = (
            <div {...this.props}>
                {this.props.children}
            </div>
        );

        legend.onAdd = function (map) {
            let div = L.DomUtil.create('div', 'leaflet-control-pos-display leaflet-bar-part leaflet-bar');
            ReactDOM.render(jsx, div);
            return div;
        };

        this.leafletElement = legend;
    }


    componentDidMount() {
        super.componentDidMount();
        this._onMouseMove();
        const {map} = this.context;
        map.on('mousemove', this._onMouseMove, this);      
    }

    _onMouseMove(e) {
        if (e) {
            var prefixAndValue = 'lat: '+ L.Util.formatNum(e.latlng.lat,5) + ' lon:' + L.Util.formatNum(e.latlng.lng,5);
            this.leafletElement._container.innerHTML = prefixAndValue;
        }
    }

    componentWillUnmount() {
        const {map} = this.context;
        map.off('mousemove', this._onMouseMove);
    }

}

