/*
 * Implement the base character rules of the Pathfinder system -
 *
 * stats, stat bonuses, etc
 *
 * In general things that "define" variables should be priority 0, and things
 * that in
 */

var stats = ["STR","DEX","CON","WIS","INT","CHA"];

addEffect({
  id: "baseDef",
  displayName: "Base Rules Stats",
  displayType: "base",
  depends: [],
  tags: ["statChange"],

  onEffect: function(entity) {
    for (var i=0; i<stats.length; i++) {
      entity[stats[i]] = 18;
    }
  }
});

addEffect({
  id: "statMods",
  displayName: "Base Rules Stat Mods",
  displayType: "base",
  depends: ["tag:statChange"],

  onEffect: function(entity) {
    for (var i=0; i<stats.length; i++) {
      entity[stats[i]+".MOD"] = Math.floor((entity[stats[i]]-10)/2);
    }
  }
})
