/*
 * Implement the base character rules of the Pathfinder system -
 *
 * stats, stat bonuses, etc
 *
 * In general things that "define" variables should be priority 0, and things
 * that in
 */

var ability_scores = ["str","dex","con","wis","int","cha"];

var bonus_types = ["alchemical", "armor", "armor_enhancement", "circumstance",
    "competence", "deflection", "dodge", "enhancement", "inherent", "insight",
    "luck", "morale", "natural_armor", "natural_armor_enhancement", "profane",
    "racial", "resistance", "sacred", "shield", "shield_enhancement", "size",
    "trait", "untyped"];

// returns accumulators for all the bonuses 
var generateBonusNames = function(stat_name) {
  return bonus_types.map(function(bonus_type) {
    return stat_name + "_" + bonus_type + "_bonus";
  });
};

var generateBonusAccums = function(entity, stat_name) {
  generateBonusNames(stat_name).forEach(function(stat_name,index) {
    // remember that some bonuses stack, others don't
    entity.vars.newAccum({
      id: stat_name,
      op: ["circumstance","dodge","racial","untyped"].indexOf(bonus_types[index]) > -1 ? '+' : 'max',
      init: 0
    });
  });
};

var sumBonuses = function(deps, stat_name) {
  return generateBonusNames(stat_name).reduce(function(accum, stat_name) {
    return accum + deps[stat_name];
  },0);
}


define.effect({
  id: "baseStats",
  displayName: "Base Stats",
  displayType: "base",

  onEffect: function(entity) {
    for (var i=0; i<ability_scores.length; i++) {
      entity.vars.new({
        id: ability_scores[i]+"_base",
        depends: [],
        onEval: function() {
          return 14;
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
    for (var i=0; i<ability_scores.length; i++) {
      // wrap the scope in a function so that callbacks reference the correct
      // value of i
      (function (score) {
        generateBonusAccums(entity,score);
        entity.vars.new({
          id: score,
          depends: [score+"_base"].concat(generateBonusNames(score)),
          onEval: function(deps) {
            return deps[score+"_base"] + sumBonuses(deps, score);
          }
        });
        entity.vars.new({
          id: score+"_mod",
          depends: [score],
          onEval: function(deps) {
            return Math.floor((deps[score]-10)/2);
          }
        });

      })(ability_scores[i]);
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
      depends: ["bab", "str_mod"],
      onEval: function(deps) {
        return deps.bab + deps.str_mod;
      }
    });

    // saves
    entity.vars.newAccum({
      id: "will_save_base",
      op: "+",
      init: 0
    });
    generateBonusAccums(entity,"will_save");
    entity.vars.new({
      id: "will_save",
      depends: ["will_save_base","wis_mod"].concat(generateBonusNames("will_save")),
      onEval: function(deps) {
        return deps.will_save_base + deps.wis_mod + sumBonuses(deps, "will_save");
      }
    })
    entity.vars.newAccum({
      id: "fort_save_base",
      op: "+",
      init: 0
    });
    generateBonusAccums(entity,"fort_save");
    entity.vars.new({
      id: "fort_save",
      depends: ["fort_save_base","con_mod"].concat(generateBonusNames("fort_save")),
      onEval: function(deps) {
        return deps.fort_save_base + deps.con_mod + sumBonuses(deps, "fort_save");
      }
    })
    entity.vars.newAccum({
      id: "ref_save_base",
      op: "+",
      init: 0
    });
    generateBonusAccums(entity,"ref_save");
    entity.vars.new({
      id: "ref_save",
      depends: ["ref_save_base","dex_mod"].concat(generateBonusNames("ref_save")),
      onEval: function(deps) {
        return deps.ref_save_base + deps.dex_mod + sumBonuses(deps, "ref_save");
      }
    })

    // testing stuff
    entity.vars.new({
      id: "fighter_lvl_bonus",
      modifies: ["fighter_lvl","will_save_insight_bonus","will_save_untyped_bonus"],
      onEval: function(deps,mods) {
        mods.will_save_insight_bonus(1);
        mods.will_save_insight_bonus(2);
        mods.will_save_untyped_bonus(1);
        mods.will_save_untyped_bonus(2);
        mods.fighter_lvl(10);
      }
    })
  }
})
