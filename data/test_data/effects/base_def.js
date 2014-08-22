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
      entity.vars[stats[i]] = 18;
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
      entity.vars[stats[i]+".MOD"] = Math.floor((entity.vars[stats[i]]-10)/2);
    }
  }
})
