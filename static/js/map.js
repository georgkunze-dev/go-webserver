map = new OpenLayers.Map("mapdiv");
map.addLayer(new OpenLayers.Layer.OSM());

var lonLat = new OpenLayers.LonLat( 10.75950444 ,51.82731744 )
        .transform(
        new OpenLayers.Projection("EPSG:4326"), // Transformation aus dem Koordinatensystem WGS 1984
        map.getProjectionObject() // in das Koordinatensystem 'Spherical Mercator Projection'
        );

var zoom=17;

var markers = new OpenLayers.Layer.Markers( "Markers" );
map.addLayer(markers);

markers.addMarker(new OpenLayers.Marker(lonLat));
map.setCenter (lonLat, zoom);