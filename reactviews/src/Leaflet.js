import React from 'react';
//import ReactDOM from 'react-dom';
import { Map, TileLayer } from 'react-leaflet';
import request from 'superagent';
import L from 'leaflet';
import MarkerClusterGroup from 'react-leaflet-markercluster';
import ZoomDisplay from './zoomDisplay';
import { MyConfig } from './MyConfig';
//import LoginByIP from './LoginByIP';
import Spinner from 'react-spinner-material';

import 'react-leaflet-markercluster/dist/styles.min.css';
import 'leaflet.markercluster/dist/MarkerCluster.css';
import 'leaflet.markercluster/dist/MarkerCluster.Default.css';

import 'leaflet/dist/leaflet.css';

// Trick for leaflet css images 
delete L.Icon.Default.prototype._getIconUrl;

L.Icon.Default.mergeOptions({
    iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
    iconUrl: require('leaflet/dist/images/marker-icon.png'),
    shadowUrl: require('leaflet/dist/images/marker-shadow.png'),
});


// that function returns Leaflet.Popup
function getLeafletPopup(url, name, ipname,loc) {
    return L.popup({minWidth: 200, closeButton: true})
        .setContent(`
    <div>
      <b>${loc}:</b><br /> 
      ${ipname} »» <a href="#/${url}ips/${name}">View IP details</a>
      </p>
    </div>
  `);
}


const stamenTonerTiles = 'http://{s}.tile.osm.org/{z}/{x}/{y}.png'; 
const stamenTonerAttr = '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors';


class Leaflet extends React.Component {

  _handleMarkerClick = (ip) => {

  }


  constructor(props) {
      super(props);
      var u = ( props.geojson === 'divgeojson' )?'div':'';
      this.state = {
          markers: null,
          url: u,
      };
  } 

  componentDidMount() {
      if ( this.props.point === 0 ) {
          const token = localStorage.getItem('token');
          request
              .get(MyConfig.API_URL + '/admin/api/v1/' + this.state.url + 'geojson') 
              .set('Authorization', `Bearer ${token}`)
              .end(function(error, response){
                  var u = this.state.url;
                  if (error) return console.error(' --' + error);
                  var ms = [];
                  response.body.features.forEach(function(p) {
                      var m = {position:[p.geometry.coordinates[1], p.geometry.coordinates[0]],
                          popup:null,options:{}};
                      m.popup = getLeafletPopup( u , p.properties.Title || '',  p.properties.IP || '', p.properties.Loc || '');
                      m.options = { IP: p.properties.IP };
                      ms.push(m);
                  });

                  this.setState({
                      markers: ms
                  });

              }.bind(this));
      }
      else { // for loc whithout marker
          var ms = [];
          var m = {position:[this.props.lat,this.props.lng],popup:'point '+this.props.point};
          ms.push(m);
          this.setState({
              markers: ms
          });
      }
  }

  render() {
      const position = [this.props.lat, this.props.lng];
      const buses = this.state.markers;

      if (!buses) {
          return (
              <div>
                  <Spinner
                      size={20}
                      spinnerColor={'#333'}
                      spinnerWidth={2}
                      visible={true} />
              </div>
          );
      }


      return (
          <Map 
              center={position} zoom={this.props.zoom} maxZoom={18} 
              style={{ height: '600px', margin: '20px' }} >

              <TileLayer
                  attribution={stamenTonerAttr}
                  url={stamenTonerTiles}
              />

              <MarkerClusterGroup
                  markers={buses}
                  showCoverageOnHover={false}
                  spiderfyDistanceMultiplier={2}
                  onMarkerClick={(marker) => this._handleMarkerClick(marker.options.IP)}
                  //wrapperOptions={{enableDefaultStyle: true}}
              />

              <ZoomDisplay />
          </Map>
      );
  }
}

Leaflet.defaultProps = {   
    lat: 46.76548,
    lng: 12.76547,
    zoom: 8,
    point: 0,
    geojson: 'geojson',
};

export default Leaflet;

