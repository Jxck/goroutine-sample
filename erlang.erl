%% receive ->
%% send (!)


% master
receive
  % send to worker
  { send, Pid } -> Pid ! { recv, self() };
  % receive from worker
  ok -> ok
end

% worker
receive 
  % receive & send
  { recv, Pid } -> Pid ! ok
end