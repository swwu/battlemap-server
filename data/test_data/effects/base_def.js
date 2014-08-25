/*
 * Implement the base character rules of the Pathfinder system -
 *
 * stats, stat bonuses, etc
 *
 * In general things that "define" variables should be priority 0, and things
 * that in
 */

var stats = ["str","dex","con","wis","int","cha"];
var stats = ["str","dex","con","wis","int","cha"];

define.effect({
  id: "baseStats",
  displayName: "Base Rules Stats",
  displayType: "base",
  tags: ["statChange.statBase"],

  onEffect: function(entity) {
    for (var i=0; i<stats.length; i++) {
      entity.vars.new({
        id: stats[i],
        dependencies: [],
        onEval: function() {
          return 7;
        }
      })
    }
  }
});

define.effect({
  id: "baseEntityRules",
  displayName: "Base Rule Calculations",
  displayType: "base",

  onEffect: function(entity) {
    // calculate stat mods
    for (var i=0; i<stats.length; i++) {
      // wrap the scope in a function so that callbacks reference the correct
      // value of i
      (function (stat) {
        entity.vars.new({
          id: stat+"_mod",
          dependencies: [stat],
          onEval: function(deps) {
            return Math.floor((deps[stat]-10)/2);
          }
        });
      })(stats[i]);
    };

    // hp
    entity.vars.newAccum({
      id: "hp",
      op: "+",
      init: 0
    });

    // calculate BAB/melee AB/ranged AB/CMB/CMD
    entity.vars.newAccum({
      id: "bab",
      op: "+",
      init: 0
    });
    entity.vars.new({
      id: "melee_ab",
      dependencies: ["bab", "str_mod"],
      onEval: function(deps) {
        return deps.bab + deps.str_mod;
      }
    });

    // saves
    entity.vars.newAccum({
      id: "will_save",
      op: "+",
      init: 0
    });
    entity.vars.newAccum({
      id: "fort_save",
      op: "+",
      init: 0
    });
    entity.vars.newAccum({
      id: "ref_save",
      op: "+",
      init: 0
    });

    // testing stuff
    entity.vars.new({
      id: "fighter_lvl_bonus",
      modifies: ["fighter_lvl"],
      onEval: function(deps,mods) {
        mods.fighter_lvl(10);
      }
    })
  }
})
