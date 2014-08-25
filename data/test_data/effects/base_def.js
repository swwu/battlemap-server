/*
 * Implement the base character rules of the Pathfinder system -
 *
 * stats, stat bonuses, etc
 *
 * In general things that "define" variables should be priority 0, and things
 * that in
 */

var stats = ["STR","DEX","CON","WIS","INT","CHA"];

define.effect({
  id: "baseStats",
  displayName: "Base Rules Stats",
  displayType: "base",
  depends: [],
  tags: ["statChange.statBase"],

  onEffect: function(entity) {
    for (var i=0; i<stats.length; i++) {
      entity.vars.new({
        id: stats[i],
        dependencies: [],
        onEval: function() {
          return 18;
        }
      })
    }
  }
});

define.effect({
  id: "statMods",
  displayName: "Base Rules Stat Mods",
  displayType: "base",
  depends: ["tag:statChange"],

  onEffect: function(entity) {
    for (var i=0; i<stats.length; i++) {
      // wrap the scope in a function so that callbacks reference the correct
      // value of i
      (function (stat) {
        entity.vars.new({
          id: stat+".MOD",
          dependencies: [stat],
          onEval: function(deps) {
            return Math.floor((deps[stat]-10)/2);
          }
        })
      })(stats[i]);
    }
  }
})
