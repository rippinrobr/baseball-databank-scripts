var async = require('async')
  , dataUtils = require('./data-utils') 
  ;

// Creates the 'Pipeline' object with the pipeline specific functions in it.
// Returns the async.compose results
var loadPitchingStats = (function() {
  var createUpdateObjects = function( pipeObj, callback ) {
    var docs = []
      , recs = pipeObj.data[ pipeObj.prevStep ]
      ;

    for(var i=0; i < recs.length; i++ ) {
      var id = recs[i].playerID;
      recs[i].season = recs[i].yearID + '_' + recs[i].teamID;

      docs.push( {query: { _id: id }, set: { $addToSet: { pitchingStats: recs[i] }}} );
    }

    pipeObj.data.push( docs );
    pipeObj.prevStep += 1;
    callback( null, pipeObj );
  };

  return async.compose(
    createUpdateObjects,
    dataUtils.createObjects,
    dataUtils.readRemoteCsv
  );

}());


var inputSrc = { path: dataUtils.baseGithubUrl + '/Pitching.csv',
              headers: 1,
          dataTypeMap: [ 'playerID', 'teamID', 'lgID', 'POS' ],
          floatColMap: [ 'BAopp', 'ERA' ] };

var outputSettings = { type: 'mongodb', 
                        url: dataUtils.mongoUrl, 
                 collection: 'players' };

var pipeObj = { input: inputSrc,
               exitFn: null,
               output: outputSettings,
                 data: [],
             prevStep: -1 };


loadPitchingStats( pipeObj, function(err, result ) {
  if ( pipeObj.output ) {
    dataUtils.updateFns[ pipeObj.output.type ]( result.data[2], pipeObj.output, function( err, r) {
      console.log('[loadPitching] Pitching Stats loaded' );
    });
  }
});
